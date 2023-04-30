package query

type SupportedLangsHandler struct {
	supportedLanguages []string
}

func NewASupportedLangsHandler(supportedLanguages []string) SupportedLangsHandler {
	return SupportedLangsHandler{supportedLanguages: supportedLanguages}
}

func (h SupportedLangsHandler) Handle() []string {
	return h.supportedLanguages
}
