import React from "react";
import { StoreCitation } from "./helpers/helper_functions"

function ArticleCard({ document }) {
//   console.log(document)
  return (
    <div className="article-card">
      <p>{document.Title}</p>{" "}
      <p>({document.Authors.map((author) => author.Name).join(", ")})</p>
      <StoreCitation document = {document}></StoreCitation>
      <br>
      </br>
      {/* TODO: Other information displayed on the cards */}
    </div>
  );
}

export default ArticleCard;
