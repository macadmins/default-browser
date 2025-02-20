package launchservices

import "fmt"

func (p *Plist) AddLSHandlers(l []LSHandler) {
	for _, handler := range l {
		fmt.Printf("Adding handler: %v\n", handler)
	}
	p.LSHandlers = append(p.LSHandlers, l...)
}

func (p *Plist) CleanHandlers() {
	var newLSHandlers []LSHandler
	for _, handler := range p.LSHandlers {
		if handler.LSHandlerContentType != "public.url" &&
			handler.LSHandlerContentType != "public.html" &&
			handler.LSHandlerContentType != "public.xhtml" &&
			handler.LSHandlerURLScheme != "https" &&
			handler.LSHandlerURLScheme != "http" {
			newLSHandlers = append(newLSHandlers, handler)
		}
	}
	p.LSHandlers = newLSHandlers
}
