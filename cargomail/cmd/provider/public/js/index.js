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

const downloadURI = (uri, name) => {
  (async () => {
    const rawResponse = await fetch(uri, {
      method: "HEAD",
      headers: {
        Accept: "application/json",
      },
    });
    if (rawResponse.status == 200) {
      const link = document.createElement("a");
      link.download = name;
      link.href = uri;
      link.click();
    } else {
      const alert = uploadForm.querySelector('div[name="alert"]')
      if (alert) alert.remove();

      uploadForm.insertAdjacentHTML(
        "beforeend",
        `<div class="alert alert-warning alert-dismissible fade show" role="alert" name="alert">
           file not found
           <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>`
      ); 
    }
  })();
};
