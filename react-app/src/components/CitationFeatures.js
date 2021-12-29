import React from "react";
let citations = [];

export const StoreCitation = (setCitations, document, e) => {
  citations.push(JSON.stringify(document.Title));
  localStorage.setItem("citations", citations);
  setCitations(citations);
};

const RemoveCitations = (setCitations, e) => {
  citations = [];
  localStorage.removeItem("citations");
  setCitations([]);
};

const ExportCitations = (setCitations, e) => {
  //implement export functionality
  RemoveCitations(setCitations, e);
};

function CitationFeatures({ setCitations }) {
  return (
    <div className="right">
      <button onClick={(e) => RemoveCitations(setCitations, e)}>
        Remove All Citations
      </button>

      <button onClick={(e) => ExportCitations(setCitations, e)}>
        Export Citations
      </button>
    </div>
  );
}

export default CitationFeatures;
