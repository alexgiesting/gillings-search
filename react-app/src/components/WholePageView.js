import React, { Component } from "react";

import { QueryForm, LoadForm, Update, ArticleOfTheDay } from "./helpers/helper_functions";

import "../App.css";
import ArticleCard from "./ArticleCard";

function WholePageView({ _ }) {
  const [results, setResults] = React.useState([]);
  
  return (
    <div className="base">
      <div className="navi">
        <span>
          <a href="https://sph.unc.edu/">
            <img
              className="logo"
              src="https://sph.unc.edu/wp-content/uploads/sites/112/2018/06/Gillings-School-of-Global-Public-Health_logo_white_h.png"
            ></img>
          </a>{" "}
        </span>
        <span></span>
        <a href="https://sph.unc.edu/resource-pages/about-the-school/">
          <span className="naviblock">About Us</span>
        </a>
        <a href="https://sph.unc.edu/students/admissions/">
          <span className="naviblock">Admissions</span>
        </a>
        <a href="https://sph.unc.edu/resource-pages/degrees-and-certificates/">
          <span className="naviblock">Degrees</span>
        </a>
        <a href="https://sph.unc.edu/students/students-home/">
          <span className="naviblock">Students</span>
        </a>
        <a href="https://sph.unc.edu/research/research/">
          {" "}
          <span className="naviblock">Research</span>
        </a>
        <a href="https://sph.unc.edu/resource-pages/sph-departments/">
          <span className="naviblock">Departments</span>
        </a>
      </div>
      <div className="header">Explore our research by department</div>
      <div className="citations">
        <div id="spacer">
          <LoadForm to="faculty" from="faculty.csv" />
          <br />
          <LoadForm to="themes" from="themes.xml" />
          <br />

          <Update
            label="Drop citations from MongoDB"
            endpoint="drop/citations"
          />
          <br />
          <Update label="Pull citations from Scopus API" endpoint="pull" />
          <br />
          <Update label="Push citations to Solr search DB" endpoint="push" />
          <br />
        </div>
      </div>

      <div className="grider">
        <div id="title">
          <h1>Gillings Search Tool</h1>
        </div>
        <div className="frame">
          <QueryForm setResults={setResults} />
          <div className="right">
            {results &&
              results.map((result, i) => (
                <ArticleCard key={i} document={result} />
              ))}
          </div>

          <div className = "articleOfTheDay" >
              {results.length == 0 ? <>
              <h1>Article of the Day</h1>
              <p><a className = "ArticleLink" href="https://projecteuclid.org/journals/statistical-science/volume-13/issue-2/Statistical-advances-in-environmental-science/10.1214/ss/1028905935.full">Statistical advances in environmental science</a></p>
              <p>(Piegorsch W.W., Smith E.P., Edwards D., Smith R.L.)</p>
              </>
              : 
              
              <></>}
          </div>

        </div>
      </div>
    </div>
  );
}

export default WholePageView;
