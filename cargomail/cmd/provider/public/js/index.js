document.addEventListener("DOMContentLoaded", function () {
  composeContent(); // in prod should be cargoesContent
});

function composeContent(e) {
  document.getElementById("composeContainer").hidden = false;
  document.getElementById("composeLink").classList.add("active");

  document.getElementById("cargoesContainer").hidden = true;
  document.getElementById("cargoesLink").classList.remove("active");

  document.getElementById("filesContainer").hidden = true;
  document.getElementById("filesLink").classList.remove("active");

  document.getElementById("composePanel").hidden = false;
  document.getElementById("cargoesPanel").hidden = true;
  document.getElementById("filesPanel").hidden = true;
}

function cargoesContent(e) {
  document.getElementById("composeContainer").hidden = true;
  document.getElementById("composeLink").classList.remove("active");

  document.getElementById("cargoesContainer").hidden = false;
  document.getElementById("cargoesLink").classList.add("active");

  document.getElementById("filesContainer").hidden = true;
  document.getElementById("filesLink").classList.remove("active");

  document.getElementById("composePanel").hidden = true;
  document.getElementById("cargoesPanel").hidden = false;
  document.getElementById("filesPanel").hidden = true;
}

function filesContent(e) {
  document.getElementById("composeContainer").hidden = true;
  document.getElementById("composeLink").classList.remove("active");

  document.getElementById("cargoesContainer").hidden = true;
  document.getElementById("cargoesLink").classList.remove("active");

  document.getElementById("filesContainer").hidden = false;
  document.getElementById("filesLink").classList.add("active");
  
  document.getElementById("composePanel").hidden = true;
  document.getElementById("cargoesPanel").hidden = true;
  document.getElementById("filesPanel").hidden = false;
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
