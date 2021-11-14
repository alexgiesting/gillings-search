import React from "react";

function ArticleCard({ document }) {
  return (
    <div className="article-card">
      <p>{document.Title}</p>{" "}
      <p>({document.Authors.map((author) => author.Name).join(", ")})</p>
      {/* TODO: Other information displayed on the cards */}
    </div>
  );
}

export default ArticleCard;
