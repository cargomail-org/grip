loginForm.onsubmit = async (e) => {
  e.preventDefault();
  const form = document.querySelector("#loginForm");

  const formData = {
    username: form.querySelector('input[name="username"]').value,
    password: form.querySelector('input[name="password"]').value,
    rememberMe: form.querySelector('input[name="rememberMe"]').checked,
  };

  console.log(formData);

  let response = await fetch("/api/v1/auth/authenticate", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  });

  const resp = await response.json();

  const alert = form.querySelector('div[name="alert"]')
  if (alert) alert.remove();

  if (resp.Err) {
    form.insertAdjacentHTML(
        "beforeend",
        `<div class="alert alert-warning alert-dismissible fade show" role="alert" name="alert">
           ${resp.Err}
           <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>`
      ); 
  } else {
    window.location.href="/"; 
  }
};
