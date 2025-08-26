/*
Copyright © 2024-present The Mochi Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package domain

type Target struct {
	Name         string  `yaml:"name"`
	Id           string  `yaml:"id"`
	TagPrefix    *string `yaml:"tagPrefix,omitempty"`
	TicketPrefix *string `yaml:"ticketPrefix,omitempty"`
}

func (t *Target) GetTagPrefix() string {
	if t.TagPrefix != nil {
		return *t.TagPrefix
	}
	return t.Id
}

func (t *Target) GetTicketPrefix() string {
	if t.TicketPrefix != nil {
		return *t.TicketPrefix
	}
	return t.Name
}
