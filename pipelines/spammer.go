package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
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

		go (func(email string) {
			defer wg.Done()
			user := GetUser(email)
			alias, hasAlias := usersAliases[user.Email]
			_, exists := users.LoadOrStore(alias, 1)

			if _, ok := users.LoadOrStore(user.Email, 1); !ok && (!hasAlias || (!exists && hasAlias)) {
				out <- user
			}
		})(email.(string))
	}

	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	users := make([]User, 0)
	wg := &sync.WaitGroup{}

	worker := func(users ...User) {
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
			go worker(users...)
			users = make([]User, 0)
		}
	}

	if len(users) != 0 {
		wg.Add(1)
		go worker(users...)
	}
	wg.Wait()
}

type spamWorker struct{}

func (w *spamWorker) doJob(msgid MsgID, out chan interface{}) {
	res, err := HasSpam(msgid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	out <- MsgData{msgid, res}
}

func CheckSpam(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	workersPool := make(chan spamWorker, HasSpamMaxAsyncRequests)

	for i := 0; i < HasSpamMaxAsyncRequests; i++ {
		workersPool <- spamWorker{}
	}

	for msgid := range in {
		id := msgid.(MsgID)
		worker := <-workersPool
		wg.Add(1)

		go (func(msgid MsgID, worker spamWorker) {
			defer wg.Done()
			worker.doJob(msgid, out)
			workersPool <- worker
		})(id, worker)

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
