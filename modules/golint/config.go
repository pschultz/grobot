package golint

const moduleConfigKey = "lint"

type Configuration struct {
	WarnCommentOrBeUnexported bool
}

var DefaultLintConfig = &Configuration{
	WarnCommentOrBeUnexported: true,
}
