document.addEventListener("DOMContentLoaded", function () {
  filesContent(); // in prod should be cargoesContent
});

function composeContent(e) {
  document.getElementById("compose-container").hidden = false;
  document.getElementById("compose-link").classList.add("active");

  document.getElementById("cargoes-container").hidden = true;
  document.getElementById("cargoes-link").classList.remove("active");

  document.getElementById("files-container").hidden = true;
  document.getElementById("files-link").classList.remove("active");
}

function cargoesContent(e) {
  document.getElementById("compose-container").hidden = true;
  document.getElementById("compose-link").classList.remove("active");

  document.getElementById("cargoes-container").hidden = false;
  document.getElementById("cargoes-link").classList.add("active");

  document.getElementById("files-container").hidden = true;
  document.getElementById("files-link").classList.remove("active");
}

function filesContent(e) {
  document.getElementById("compose-container").hidden = true;
  document.getElementById("compose-link").classList.remove("active");

  document.getElementById("cargoes-container").hidden = true;
  document.getElementById("cargoes-link").classList.remove("active");

  document.getElementById("files-container").hidden = false;
  document.getElementById("files-link").classList.add("active");
}
