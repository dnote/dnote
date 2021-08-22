function toggleAccountDropdown() {
  document;
}

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

console.log('aaaa');

var dropdownTriggerEls = document.getElementsByClassName('dropdown-trigger');

for (var i = 0; i < dropdownTriggerEls.length; i++) {
  var dropdownTriggerEl = dropdownTriggerEls[i];

  dropdownTriggerEl.addEventListener('click', function (e) {
    var el = getNextSibling(e.target, '.dropdown-content');

    el.classList.toggle('show');
  });
}


// TODO: how to close dropdown?
// if click target contains trigger or content, noop. otherwise close.

// window.onclick = function (e) {
//   // Close dropdown on click outside the dropdown content
//   if (
//     !e.target.matches('.dropdown-trigger') &&
//     !e.target.matches('.dropdown-content')
//   ) {
//     var dropdowns = document.getElementsByClassName('dropdown-content');
//     for (var i = 0; i < dropdowns.length; i++) {
//       var openDropdown = dropdowns[i];
//       if (openDropdown.classList.contains('show')) {
//         openDropdown.classList.remove('show');
//       }
//     }
//   }
// };
