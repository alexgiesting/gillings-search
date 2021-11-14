import React, { Component } from "react";

import { QueryForm, LoadForm, Update } from "./helpers/helper_functions";

import "../App.css";
import ArticleCard from "./ArticleCard";

// export default class WholePageView extends Component {}

function WholePageView({ _ }) {
  const [results, setResults] = React.useState([]);
  return (
    <div className="base">
      <LoadForm to="faculty" from="faculty.csv" />
      <LoadForm to="themes" from="themes.xml" />
      <br />

      <Update label="Drop citations from MongoDB" endpoint="drop/citations" />
      <Update label="Pull citations from Scopus API" endpoint="pull" />
      <Update label="Push citations to Solr search DB" endpoint="push" />
      <br />

      <h1>Gillings Search Tool</h1>
      <div className="frame">
        <QueryForm setResults={setResults} />
        <div className="right">
          {results &&
            results.map((result, i) => (
              <ArticleCard key={i} document={result} />
            ))}
        </div>
      </div>
    </div>
  );
}

export default WholePageView;
