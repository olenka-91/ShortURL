package service

type urlStoreTranslation interface {
	ShortURL(longURL []byte) (string, error)
	LongURL(shortURL string) (string, error)
}

type Service struct {
	urlStoreTranslation
}

func NewService( /*r *repository.Repository*/ ) *Service {
	return &Service{urlStoreTranslation: NewUrlStoreTranslationService( /*r.Song*/ )}
}
