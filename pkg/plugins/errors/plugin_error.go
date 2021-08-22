// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"encoding/json"
	"fmt"

	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
)

type PluginError struct {
	code    codes.Code `json:"-"`
	Code    string     `json:"code"`
	Msg     string     `json:"msg,omitempty"`
	Cause   string     `json:"cause,omitempty"`
	Service string     `json:"service,omitempty"`
	Plugin  string     `json:"plugin,omitempty"`
	Args    string     `json:"args,omitempty"`
}

func (p *PluginError) Error() string {
	p.Code = p.code.String()
	data, _ := json.Marshal(p)
	return string(data)
}

// Code - returns a nitric api error code from an error or Unknown if the error was not a nitric api error
func Code(e error) codes.Code {
	if pe, ok := e.(*PluginError); ok {
		return pe.code
	}

	return codes.Unknown
}

// ErrorsWithScope - Returns a new reusable error factory with the given scope
func ErrorsWithScope(s string, ctx ...interface{}) func(c codes.Code, msg string, cause error) error {
	return func(c codes.Code, msg string, cause error) error {
		pe := &PluginError{
			code:   c,
			Msg:    msg,
			Plugin: s,
		}
		if ctx != nil {
			pe.Args = fmt.Sprintf("%v", ctx)
		}
		if cause != nil {
			pe.Cause = fmt.Sprintf("%v", cause)
		}
		return pe
	}
}
