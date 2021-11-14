import React from 'react';

function ArticleCard(props) {
  return (
    <div className="article-card">
      <p>{props.Title}</p> {" "}
      <p>({props.Authors.map((author) => author.Name).join(", ")})</p>
      {/* TODO: Other information displayed on the cards */}
    </div>
  )
}

export default ArticleCard;