import React from "react";
import { StoreCitation } from "./CitationFeatures";

function ArticleCard({ setCitations, document }) {
  return (
    <div className="article-card">
      <p>{document.Title}</p>{" "}
      <p>({document.Authors.map((author) => author.Name).join(", ")})</p>
      <button onClick={(e) => StoreCitation(setCitations, document, e)}>
        Save
      </button>
      {/* TODO: Other information displayed on the cards */}
    </div>
  );
}

export default ArticleCard;
