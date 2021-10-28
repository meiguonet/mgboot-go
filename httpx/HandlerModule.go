package httpx

type HandlerModule interface {
	GetHandlerEntries() []*HandlerEntry
}
