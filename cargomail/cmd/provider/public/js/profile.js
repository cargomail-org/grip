const profileForm = document.getElementById("profileForm");

profileForm.onsubmit = async (e) => {
  e.preventDefault();

  const form = e.currentTarget;

  const formData = {
    firstname: form.querySelector('input[name="firstname"]').value,
    lastname: form.querySelector('input[name="lastname"]').value,
  };

  const response = await api(form.id, 200, "/api/v1/user/profile", {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  });

  if (response === false) {
    return;
  }
};
