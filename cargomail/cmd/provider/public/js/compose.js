import DataTable from "datatables.net";

const toInput = document.getElementById("toInput");
const subjectInput = document.getElementById("subjectInput");
const messageText = document.getElementById("messageText");

const subjectHeadings = document.getElementsByClassName("subject-heading");
let bouncerTimeout = null;

const bouncer = (e) => {
  clearTimeout(bouncerTimeout);
  bouncerTimeout = setTimeout(() => {
    [...subjectHeadings].forEach((heading) => {
      heading.textContent = subjectInput.value;
    });
    //TODO save data to server
  }, 2000);
};

toInput.addEventListener("keyup", (event) => bouncer(event));
subjectInput.addEventListener("keyup", (event) => bouncer(event));
messageText.addEventListener("keyup", (event) => bouncer(event));

const composeTable = new DataTable("#composeTable", {
  paging: true,
  responsive: {
    details: false,
  },
});
