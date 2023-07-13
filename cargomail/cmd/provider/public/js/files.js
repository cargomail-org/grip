import DataTable from "datatables.net";

import "datatables.net-bs5";
import "datatables.net-select";
import "datatables.net-select-bs5";
import "datatables.net-responsive";
import "datatables.net-responsive-bs5";

const table = new DataTable("#files-table", {
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
    { data: "name" },
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
});
