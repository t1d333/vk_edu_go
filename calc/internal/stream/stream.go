package stream


type Stream  interface {
  PeekChar() (rune, error)
  ReadChar() (rune, error)
  Empty() bool
}
