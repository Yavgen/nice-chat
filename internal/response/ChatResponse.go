package response

type ChatResponse interface {
	ToJson() JsonResponse
}
