import DataTable from "datatables.net";

// import $ from 'jquery';

import "datatables.net-bs5";
import "datatables.net-select";
import "datatables.net-select-bs5";
import "datatables.net-buttons";
import "datatables.net-buttons-bs5";
import "datatables.net-responsive";
import "datatables.net-responsive-bs5";

let selectedIds = [];

const confirmDialog = new bootstrap.Modal(
  document.querySelector("#confirmDialog")
);

const filesTable = new DataTable("#filesTable", {
  paging: true,
  responsive: {
    details: false,
  },
  ajax: {
    url: "/api/v1/resources",
    dataSrc: "",
  },
  columns: [
    { data: "id", visible: false, searchable: false },
    { data: null, visible: true, orderable: false, width: "15px" },
    {
      data: "name",
      render: (data, type, full, meta) => {
        const link = "/api/v1/resource/";
        return `<a href="${link}${full.uuid}">${data}</a>`;
      },
    },
    { data: "size", searchable: false },
  ],
  columnDefs: [
    {
      orderable: false,
      className: "select-checkbox",
      targets: 1,
      data: null,
      defaultContent: "",
    },
  ],
  select: {
    style: "multi",
    selector: "td:first-child",
    info: true,
  },
  order: [[0, "desc"]],
  dom: "Bfrtip",
  language: {
    buttons: {
      pageLength: 'Show %d'
    }
  },
  lengthMenu: [
    [5, 10, 25],
    ["5 rows", "10 rows", "25 rows"],
  ],
  buttons: [
    "pageLength",
    {
      text: "Reload",
      action: function () {
        filesTable.ajax.reload();
      },
    },
    {
      text: "Delete",
      className: 'files-delete',
      enabled: false,
      action: function () {
        selectedIds = [];

        const selectedData = filesTable
          .rows(".selected")
          .data()
          .map((obj) => obj.id);
        if (selectedData.length > 0) {
          confirmDialog.show();
          for (let i = 0; i < selectedData.length; i++) {
            selectedIds.push(parseInt(selectedData[i], 10));
          }
        }
      },
    },
  ],
});

filesTable.on('select.dt deselect.dt', () => {
  filesTable.buttons( ['.files-delete'] ).enable(
    filesTable.rows( { selected: true } ).indexes().length === 0 ? false : true
  )
});

export const deleteItems = (e) => {
  e.preventDefault();

  confirmDialog.hide();

  console.log(selectedIds);

  fetch("api/v1/resources/delete", {
    method: "post",
    body: JSON.stringify(selectedIds),
  })
    .then((response) => {
      if (response.ok) {
        filesTable.rows(".selected").remove().draw();
        filesTable.buttons( ['.files-delete'] ).enable(false);
        console.log("Successfully deleted file(s)");
        if (!Object.keys(response).length) {
          console.log("no return data found");
          return;
        }
        return response && response.json();
      }
      throw new Error("something went wrong");
    })
    .then((data) => {
      if (data) {
        console.log("File(s) deleted:", data);
      }
    })
    .catch((error) => {
      console.log("Error while deleting file(s):", error);
    });
};

export const clearUpload = (e) => {
  e.preventDefault();
  document.getElementById("inputUpload").value = "";
};
