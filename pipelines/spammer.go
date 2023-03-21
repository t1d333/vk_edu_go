package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"sync/atomic"
)

func RunPipeline(cmds ...cmd) {
	channels := make([]chan interface{}, len(cmds)+1)
	wg := &sync.WaitGroup{}

	for i := range channels {
		channels[i] = make(chan interface{})
	}

	for i := 0; i < len(cmds); i++ {
		wg.Add(1)
		go (func(ch1, ch2 chan interface{}, f cmd) {
			f(ch1, ch2)
			wg.Done()
			close(ch2)
		})(channels[i], channels[i+1], cmds[i])
	}
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	users := sync.Map{}
	for email := range in {
		wg.Add(1)
		go (func(email string, users *sync.Map, out chan interface{}) {
			defer wg.Done()
			user := GetUser(email)
			alias, hasAlias := usersAliases[user.Email]
			_, exists := users.Load(alias)
			if _, ok := users.Load(user.Email); !ok && hasAlias && !exists {
				out <- user
				users.Store(alias, 1)
			} else if !ok {
				out <- user
				users.Store(user.Email, 1)
			}
		})(email.(string), &users, out)
	}

	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	users := make([]User, 0)
	wg := &sync.WaitGroup{}
	worker := func(out chan interface{}, users ...User) {
		defer wg.Done()
		msgs, err := GetMessages(users...)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            return
        }
		for _, msg := range msgs {
			out <- MsgID(msg)
		}
	}

	for user := range in {
		users = append(users, user.(User))
		if len(users) == GetMessagesMaxUsersBatch {
			wg.Add(1)
			go worker(out, users...)
			users = make([]User, 0)
		}
	}

	if len(users) != 0 {
		wg.Add(1)
		go worker(out, users...)
	}
	wg.Wait()
}

func CheckSpam(in, out chan interface{}) {
	counter := int32(0)

	wg := sync.WaitGroup{}
	worker := func(msgid MsgID, counter *int32) {
		defer wg.Done()
		defer atomic.AddInt32(counter, -1)
		res, err := HasSpam(msgid)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            return
        }
		out <- MsgData{msgid, res}
	}
	for msgid := range in {
		id := msgid.(MsgID)

		wg.Add(1)
		if atomic.LoadInt32(&counter) < int32(HasSpamMaxAsyncRequests) {
			atomic.AddInt32(&counter, 1)
			go worker(id, &counter)
		} else {
			for atomic.LoadInt32(&counter) >= int32(HasSpamMaxAsyncRequests) {
			}
			atomic.AddInt32(&counter, 1)
			go worker(id, &counter)
		}
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	result := make([]MsgData, 0)
	for data := range in {
		result = append(result, data.(MsgData))
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].HasSpam && !result[j].HasSpam {
			return true
		}

		if result[i].HasSpam == result[j].HasSpam {
			return result[i].ID < result[j].ID
		}

		return false
	})

	for _, data := range result {
		out <- fmt.Sprintf("%v %v", data.HasSpam, data.ID)
	}
}
