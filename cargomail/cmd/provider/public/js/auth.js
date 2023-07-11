const registerForm = document.getElementById('registerForm');
const loginForm = document.getElementById('loginForm');

if (registerForm) {
  registerForm.onsubmit = async (e) => {
    e.preventDefault();
    const form = document.querySelector("#registerForm");
  
    const formData = {
      username: form.querySelector('input[name="username"]').value,
      password: form.querySelector('input[name="password"]').value,
    };
  
    console.log(formData);
  
    const confirmation = form.querySelector('input[name="confirmation"]').value;
  
    const alert = form.querySelector('div[name="alert"]')
    if (alert) alert.remove();
  
    if (formData.password != confirmation) {
      form.insertAdjacentHTML(
        "beforeend",
        `<div class="alert alert-warning alert-dismissible fade show" role="alert" name="alert">
          Passwords do NOT match
           <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>`
      );
      return
    }
  
    let response = await fetch("/api/v1/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    });
  
    const resp = await response.json();
  
    if (resp.Err) {
      form.insertAdjacentHTML(
          "beforeend",
          `<div class="alert alert-warning alert-dismissible fade show" role="alert" name="alert">
             ${resp.Err.charAt(0).toUpperCase() + resp.Err.slice(1)}
             <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
          </div>`
        ); 
    } else {
      window.location.href="/login"; 
    }
  };
}

if (loginForm) {
  loginForm.onsubmit = async (e) => {
    e.preventDefault();
    const form = document.querySelector("#loginForm");
  
    const formData = {
      username: form.querySelector('input[name="username"]').value,
      password: form.querySelector('input[name="password"]').value,
      rememberMe: form.querySelector('input[name="rememberMe"]').checked,
    };
  
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
          ${resp.Err.charAt(0).toUpperCase() + resp.Err.slice(1)}
          <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
          </div>`
        ); 
    } else {
      window.location.href="/"; 
    }
  };
}


