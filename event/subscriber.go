package event

const ()

type Subscriber interface {
	Subscribe(event string) <-chan interface{}
}

type BlockSubscriber struct {
}
