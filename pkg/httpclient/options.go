package httpclient

type Option func(client *Client)

func WithBaseURL(addr string) Option {
	return func(s *Client) {
		s.SetBaseURL(addr)
	}
}
