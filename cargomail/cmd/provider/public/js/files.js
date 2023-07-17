import DataTable from "datatables.net";

// import $ from 'jquery';

import "datatables.net-bs5";
import "datatables.net-select";
import "datatables.net-select-bs5";
import "datatables.net-buttons";
import "datatables.net-buttons-bs5";
import "datatables.net-responsive";
import "datatables.net-responsive-bs5";

export const startOrResumeUpload = (upload) => {
  upload.findPreviousUploads().then(function (previousUploads) {
    if (previousUploads.length) {
      upload.resumeFromPreviousUpload(previousUploads[0]);
    }
    upload.start();
  });
};

export const createTusUploadInstance = (url, file) => {
  const upload = new tus.Upload(file, {
    endpoint: url,
    retryDelays: [0, 3000, 5000],
    chunkSize: 8192, // dev - upload slowdown
    enableChecksum: true,
    metadata: {
      filename: file.name,
      filetype: file.type,
    },
    onError: (error) => console.log("failed because: " + error),
  });
  return upload;
};

let selectedUuids = [];

const confirmDialog = new bootstrap.Modal(
  document.querySelector("#confirmDialog")
);

const uploadForm = document.getElementById("uploadForm");

uploadForm.onsubmit = async (e) => {
  e.preventDefault();
  const form = e.currentTarget;
  // using tus
  // const url = uploadForm.action + "/tus/upload";
  // using FormData
  const url = uploadForm.action + "/upload";

  const formData = new FormData(uploadForm);

  for (const pair of formData.entries()) {
    if (typeof pair[1] == "object") {
      const file = pair[1];
      // using tus
      // const upload = createTusUploadInstance(url, file);
      // upload.start();

      // using FormData
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
        const content = await rawResponse.json();

        filesTable.row.add(content);
        filesTable.draw();

        // let newdata = Array.from(filesTable.data());
        // newdata.unshift(content);
        // filesTable.clear();
        // for (const row of newdata) {
        //   filesTable.row.add(row);
        // }
        // filesTable.draw();
      })();

      /*const response = await fetch(url, {
        method: "POST",
        body: singleFileFormData,
      })
        .then((response) => {
          if (response.ok) {
            console.log("Successfully uploaded file");
            // if (!Object.keys(response).length) {
            //   console.log("no return data found");
            //   return;
            // }
            // const uploadedFiles = response.json();

            // let newdata = Array.from(filesTable.data());
            // newdata.unshift(["my", "new", "row"]);
            // filesTable.clear();
            // for (const row of newdata) {
            //   filesTable.row.add(row);
            // }
            // filesTable.draw();
          }
          // throw new Error("something went wrong");
        })
        .then((data) => {
          if (data) {
            console.log("File uploaded:", data);
          }
        })
        .catch((error) => {
          console.log("Error while uploading file:", error);
        });*/
    }
  }
};

/* uploadForm.onsubmit = async (e) => {
  e.preventDefault();
  const form = e.currentTarget;
  const url = uploadForm.action;

  const formData = new FormData(uploadForm);
  const response = await fetch(url, {
    method: "POST",
    body: formData,
  })
    .then((response) => {
      if (response.ok) {
        console.log("Successfully uploaded file(s)");
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
        console.log("File(s) uploaded:", data);
      }
    })
    .catch((error) => {
      console.log("Error while uploading file(s):", error);
    });
}; */

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
        return `<a href="${link}${full.uuid}">${data}</a>`;
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
  e.preventDefault();

  confirmDialog.hide();

  console.log(selectedUuids);

  fetch("api/v1/files/delete", {
    method: "post",
    Accept: "application/json",
    "Content-Type": "application/json",
    body: JSON.stringify(selectedUuids),
  })
    .then((response) => {
      if (response.ok) {
        filesTable.rows(".selected").remove().draw();
        filesTable.buttons([".files-delete"]).enable(false);
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

export const inputUploadChanged = (e) => {
  e.preventDefault();
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
  e.preventDefault();
  document.getElementById("uploadButton").classList.add("disabled");
  document.getElementById("inputUpload").value = "";
  document.getElementById("clearButton").classList.add("disabled");
};
