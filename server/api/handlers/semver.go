/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package handlers

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type semver struct {
	Major int
	Minor int
	Patch int
}

func parseSemver(version string) (semver, error) {
	re := regexp.MustCompile(`(\d*)\.(\d*)\.(\d*)`)
	match := re.FindStringSubmatch(version)

	if len(match) != 4 {
		return semver{}, errors.Errorf("invalid semver %s", version)
	}

	major, err := strconv.Atoi(match[1])
	if err != nil {
		return semver{}, errors.Wrap(err, "converting major version to int")
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return semver{}, errors.Wrap(err, "converting minor version to int")
	}
	patch, err := strconv.Atoi(match[3])
	if err != nil {
		return semver{}, errors.Wrap(err, "converting patch version to int")
	}

	ret := semver{
		Major: major,
		Minor: minor,
		Patch: patch,
	}

	return ret, nil
}
