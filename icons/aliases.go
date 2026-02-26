package icons

import "github.com/a-h/templ"

// Aliases for backward compatibility with the previous manual icon set
// These map old icon names to their new Lucide equivalents

// Home creates a home icon (alias for House)
func Home(opts ...Option) templ.Component {
	return House(opts...)
}

// AlertCircle creates an alert/warning icon with a circle (alias for CircleAlert)
func AlertCircle(opts ...Option) templ.Component {
	return CircleAlert(opts...)
}

// CheckCircle creates a success/check icon with a circle (alias for CircleCheck)
func CheckCircle(opts ...Option) templ.Component {
	return CircleCheck(opts...)
}

// XCircle creates an error/close icon with a circle (alias for CircleX)
func XCircle(opts ...Option) templ.Component {
	return CircleX(opts...)
}
