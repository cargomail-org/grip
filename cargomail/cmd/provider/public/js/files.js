import DataTable from "datatables.net";

// import $ from 'jquery';

import "datatables.net-bs5";
import "datatables.net-select";
import "datatables.net-select-bs5";
import "datatables.net-buttons";
import "datatables.net-buttons-bs5";
import "datatables.net-responsive";
import "datatables.net-responsive-bs5";

let selectedUuids = [];

const confirmDialog = new bootstrap.Modal(
  document.querySelector("#confirmDialog")
);

const uploadForm = document.getElementById("uploadForm");

uploadForm.onsubmit = async (e) => {
  e?.preventDefault();
  const form = e.currentTarget;
  const url = uploadForm.action;

  const formData = new FormData(uploadForm);

  const entries = [...formData.entries()];
  entries.forEach(function (entry, index) {
    if (entry[1] instanceof File) {
      const file = entry[1];
      const singleFileFormData = new FormData();
      singleFileFormData.append("files", file);

      (async () => {
        const rawResponse = await fetch(url, {
          method: "POST",
          headers: {
            Accept: "application/json",
          },
          body: singleFileFormData,
        });
        if (rawResponse.ok) {
          const content = await rawResponse.json();

          filesTable.row.add(content);
          filesTable.draw();
        }
      })();
    }
  });
  clearUpload();
};

const filesTable = new DataTable("#filesTable", {
  paging: true,
  responsive: {
    details: false,
  },
  ajax: {
    url: "/api/v1/files",
    dataSrc: "",
  },
  columns: [
    { data: "uuid", visible: false, searchable: false },
    { data: null, visible: true, orderable: false, width: "15px" },
    {
      data: "name",
      render: (data, type, full, meta) => {
        const link = "/api/v1/files/";
        // return `<a href="${link}${full.uuid}" target="_blank">${data}</a>`;
        return `<a href="javascript:;" onclick="downloadURI('${link}${full.uuid}', '${data}');">${data}</a>`;
      },
    },
    { data: "size", searchable: false },
    { data: "created_at", searchable: true },
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
  order: [[4, "desc"]],
  dom: "Bfrtip",
  language: {
    buttons: {
      pageLength: "Show %d",
    },
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
        filesTable.buttons([".files-delete"]).enable(false);
      },
    },
    {
      text: "Delete",
      className: "files-delete",
      enabled: false,
      action: function () {
        selectedUuids = [];

        const selectedData = filesTable
          .rows(".selected")
          .data()
          .map((obj) => obj.uuid);
        if (selectedData.length > 0) {
          confirmDialog.show();
          for (let i = 0; i < selectedData.length; i++) {
            selectedUuids.push(selectedData[i]);
          }
        }
      },
    },
  ],
});

filesTable.on("select.dt deselect.dt", () => {
  filesTable
    .buttons([".files-delete"])
    .enable(
      filesTable.rows({ selected: true }).indexes().length === 0 ? false : true
    );
});

export const deleteItems = (e) => {
  e?.preventDefault();

  confirmDialog.hide();

  console.log(selectedUuids);

  (async () => {
    const rawResponse = await fetch("api/v1/files/delete", {
      method: "DELETE",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(selectedUuids),
    });
    if (rawResponse.ok) {
      const content = await rawResponse.json();
      if (content?.status == "OK") {
        filesTable.rows(".selected").remove().draw();
        filesTable.buttons([".files-delete"]).enable(false);
        console.log("Successfully deleted file(s)");
      }
    }
  })();
};

export const inputUploadChanged = (e) => {
  e?.preventDefault();
  const files = e.target.files;
  if (files.length && files.length > 0) {
    document.getElementById("uploadButton").classList.remove("disabled");
    document.getElementById("clearButton").classList.remove("disabled");
  } else {
    document.getElementById("uploadButton").classList.add("disabled");
    document.getElementById("clearButton").classList.add("disabled");
  }
};

export const clearUpload = (e) => {
  e?.preventDefault();
  document.getElementById("uploadButton").classList.add("disabled");
  uploadForm.reset();
  document.getElementById("clearButton").classList.add("disabled");
};
