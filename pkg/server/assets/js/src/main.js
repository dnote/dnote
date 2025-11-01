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
