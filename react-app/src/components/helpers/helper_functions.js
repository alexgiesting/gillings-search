import React, { useState } from "react";

export function QueryForm({ setResults }) {
  return (
    <form
      style={{ width: "600px" }}
      onSubmit={async function (event) {
        event.preventDefault();
        const request = {};
        [...event.target]
          .filter((input) => input.value !== "" && input.name !== "")
          .forEach((input) => {
            request[input.name] = input.value;
          });
        const results = await fetch(
          `/solr/citations/select?${makeQuery(request)}`,
          { method: "GET" }
        ).then((response) => response.json()); // TODO error handling
        // console.log(results);

        let final;
        {
          const ids = { id: results.response.docs.map((doc) => doc.id) };
          const q = `q=${encodeURIComponent(JSON.stringify(ids))}`;
          final = await fetch(`/query?${q}`, { method: "GET" }).then(
            (response) => response.json()
          );
        }
        //console.log(final);
        setResults(final);
      }}
    >
      {/* Main search bar */}
      <div style={{ display: "flex", width: "100%", marginTop: 20 }}>
        <input
          name="keyword"
          type="text"
          style={{
            flexGrow: 1,
            backgroundColor: "#e0e0e0",
            borderStyle: "none",
            height: 40,
            borderRadius: 5,
          }}
        />
        <input
          type="submit"
          value="Search"
          style={{
            marginLeft: "1em",
            backgroundColor: "#aed3fc",
            borderStyle: "none",
            borderRadius: 5,
            color: "white",
          }}
        />
      </div>
      <br />

      {/* <h4>Filters</h4> */}
      <div style={{ textAlign: "center" }}>
        {/* <QueryText name="strengths" label="Research Strength(s)" />
        <QueryText name="department" label="Department" /> */}
        <span style={{ display: "inline-block" }}>
          <QueryText name="title" label="Title" />
        </span>
        <span style={{ display: "inline-block" }}>
          <QueryText name="author" label="Author(s)" />
        </span>
        <span style={{ display: "inline-block" }}>
          <div
            style={{
              position: "relative",
              width: "100%",
              marginBottom: "0.5em",
            }}
          >
            <label>
              Year (Range)
              <input
                name="start"
                type="number"
                min="1900"
                max={new Date(Date.now()).getFullYear()}
                style={{
                  right: 0,
                  backgroundColor: "#e0e0e0",
                  borderStyle: "none",
                  borderRadius: 5,
                  height: 30,
                }}
              />
              <input
                name="end"
                type="number"
                min="1900"
                max={new Date(Date.now()).getFullYear()}
                style={{
                  right: 0,
                  top: "100%",
                  backgroundColor: "#e0e0e0",
                  borderStyle: "none",
                  borderRadius: 5,
                  height: 30,
                }}
              />
            </label>
            <br />
          </div>
        </span>
      </div>
      <br />
    </form>
  );
}

export function makeQuery(request) {
  let q = [];

  function makeField(fieldName, getValues = (x) => [x], useLabel = true) {
    if (fieldName in request) {
      const prefix = useLabel ? fieldName + ":" : "";
      for (const value of getValues(request[fieldName])) {
        q.push(prefix + value);
      }
    }
  }

  function splitter(text) {
    return text
      .replace(/[\\+\-&|!(){}[\]^~?:;]+/g, "")
      .replace(",", " ")
      .split(/[,\s]+/); // TODO this isn't right?
  }

  makeField("keyword", splitter, false);
  makeField("author", splitter);
  makeField("title", splitter);

  if ("start" in request || "end" in request) {
    const suffix = "-01-01T00:00:00Z";
    const start = "start" in request ? request["start"] + suffix : "*";
    const end = "end" in request ? request["end"] + suffix : "NOW";
    request["date"] = `[${start} TO ${end}]`;
  }
  makeField("date");

  return `q.op=AND&q=${encodeURIComponent(q.join(" "))}&fl=id`;
}

export function QueryText({ name, label }) {
  return (
    <div style={{ width: "100%", marginBottom: "0.5em" }}>
      <label>
        {label}
        <input
          name={name}
          type="text"
          style={{
            right: 0,
            backgroundColor: "#e0e0e0",
            borderStyle: "none",
            borderRadius: 5,
            height: 30,
          }}
        />
      </label>
      <br />
    </div>
  );
}

export function LoadForm({ to, from }) {
  return (
    <form
      method="post"
      action={`/update/load/${to}`}
      encType="multipart/form-data"
    >
      <label>
        upload {from}:{" "}
        <input name="file" type="file" style={{ marginLeft: 60 }} />
      </label>
      <br />
      <label>
        <input
          name="key"
          type="text"
          placeholder="password"
          style={{
            backgroundColor: "#e0e0e0",
            borderStyle: "none",
            borderRadius: 5,
            height: 30,
          }}
        />
      </label>
      <input
        type="submit"
        value="Upload"
        style={{
          backgroundColor: "#13294B",
          borderStyle: "none",
          borderRadius: 5,
          color: "white",
          height: 30,
        }}
      />
    </form>
  );
}

export function Update({ endpoint, label }) {
  return (
    <form action={`/update/${endpoint}`}>
      <Password />
      <input
        type="submit"
        value={label}
        style={{
          marginLeft: "1em",
          backgroundColor: "#13294B",
          borderStyle: "none",
          borderRadius: 5,
          color: "white",
          height: 30,
        }}
      />
    </form>
  );
}

export function Password({ _ }) {
  return (
    <label>
      <input
        name="key"
        type="text"
        placeholder="(password)"
        style={{
          backgroundColor: "#e0e0e0",
          borderStyle: "none",
          borderRadius: 5,
          height: 30,
        }}
      />
    </label>
  );
}

// export function StoreCitation({ document }) {
//   const [ anchorEl, setAnchorEl ] = useState(null);

//   const handleSave = (event) => {
//         setAnchorEl(event.currentTarget);
//     };

//   return (
//     <button type="button" onClick= {handleSave}>Save</button>
//   );

// }
