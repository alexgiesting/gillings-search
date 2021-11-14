import React, { Component } from "react";

import { Result, QueryForm, LoadForm, Update  } from "./helpers/helper_functions";

import "../App.css"



// export default class WholePageView extends Component {} 

function WholePageView({ _ }) {
  const [results, setResults] = React.useState([]);
  return (
    <div className="base">
      <LoadForm to="faculty" from="faculty.csv" />
      <LoadForm to="themes" from="themes.xml" />
      <br />

      <Update label="Drop citations" endpoint="drop/citations" />
      <Update label="Pull citations" endpoint="pull" />
      <br />

      <h1>Gillings Search Tool</h1>
      <div className="frame">
        <QueryForm setResults={setResults} />
        <div className="right">
          {results &&
            results.map((result, i) => <Result key={i} result={result} />)}
        </div>
      </div>
    </div>
  );
}


export default WholePageView;