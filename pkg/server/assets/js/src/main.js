/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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

var getNextSibling = function (el, selector) {
  var sibling = el.nextElementSibling;

  if (!selector) {
    return sibling;
  }

  while (sibling) {
    if (sibling.matches(selector)) return sibling;
    sibling = sibling.nextElementSibling;
  }
};

var dropdownTriggerEls = document.getElementsByClassName('dropdown-trigger');

for (var i = 0; i < dropdownTriggerEls.length; i++) {
  var dropdownTriggerEl = dropdownTriggerEls[i];

  dropdownTriggerEl.addEventListener('click', function (e) {
    var el = getNextSibling(e.target, '.dropdown-content');

    el.classList.toggle('show');
  });
}

// Dropdown closer
window.onclick = function (e) {
  // Close dropdown on click outside the dropdown content or trigger
  function shouldClose(target) {
    var dropdownContentEls = document.getElementsByClassName(
      'dropdown-content'
    );

    for (let i = 0; i < dropdownContentEls.length; ++i) {
      var el = dropdownContentEls[i];
      if (el.contains(target)) {
        return false;
      }
    }
    for (let i = 0; i < dropdownTriggerEls.length; ++i) {
      var el = dropdownTriggerEls[i];
      if (el.contains(target)) {
        return false;
      }
    }

    return true;
  }

  if (shouldClose(e.target)) {
    var dropdowns = document.getElementsByClassName('dropdown-content');
    for (var i = 0; i < dropdowns.length; i++) {
      var openDropdown = dropdowns[i];
      if (openDropdown.classList.contains('show')) {
        openDropdown.classList.remove('show');
      }
    }
  }
};
