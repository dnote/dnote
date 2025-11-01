/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package views

import (
	"fmt"
	"time"
)

type timeDiff struct {
	text  string
	tense string
}

func pluralize(singular string, count int) string {
	var noun string
	if count == 1 {
		noun = singular
	} else {
		noun = singular + "s"
	}

	return noun
}

func abs(num int64) int64 {
	if num < 0 {
		return -num
	}

	return num
}

var (
	DAY  = 24 * time.Hour.Milliseconds()
	WEEK = 7 * DAY
)

func getTimeDiffText(interval int64, noun string) string {
	return fmt.Sprintf("%d %s", interval, pluralize(noun, int(interval)))
}

func relativeTimeDiff(t1, t2 time.Time) timeDiff {
	diff := t1.Sub(t2)
	ts := abs(diff.Milliseconds())

	var tense string
	if diff > 0 {
		tense = "past"
	} else {
		tense = "future"
	}

	interval := ts / (52 * WEEK)
	if interval >= 1 {
		return timeDiff{
			text:  getTimeDiffText(interval, "year"),
			tense: tense,
		}
	}

	interval = ts / (4 * WEEK)
	if interval >= 1 {
		return timeDiff{
			text:  getTimeDiffText(interval, "month"),
			tense: tense,
		}
	}

	interval = ts / WEEK
	if interval >= 1 {
		return timeDiff{
			text:  getTimeDiffText(interval, "week"),
			tense: tense,
		}
	}

	interval = ts / DAY
	if interval >= 1 {
		return timeDiff{
			text:  getTimeDiffText(interval, "day"),
			tense: tense,
		}
	}

	interval = ts / time.Hour.Milliseconds()
	if interval >= 1 {
		return timeDiff{
			text:  getTimeDiffText(interval, "hour"),
			tense: tense,
		}
	}

	interval = ts / time.Minute.Milliseconds()
	if interval >= 1 {
		return timeDiff{
			text:  getTimeDiffText(interval, "minute"),
			tense: tense,
		}
	}

	return timeDiff{
		text: "Just now",
	}
}
