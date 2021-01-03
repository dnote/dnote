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
