// import logo from "./logo.svg";
// import "./App.css";

// function App() {
//   return (
//     <div className="App">
//       <header className="App-header">
//         <img src={logo} className="App-logo" alt="logo" />
//         <p>
//           Edit <code>src/App.js</code> and save to reload.
//         </p>
//         <a
//           className="App-link"
//           href="https://reactjs.org"
//           target="_blank"
//           rel="noopener noreferrer"
//         >
//           Learn React
//         </a>
//       </header>
//     </div>
//   );
// }

function App() {
  return <AppHeader header="Gillings Search Tool" />;
}

function AppHeader({ header }) {
  return (
    <div>
      <h1>{header}</h1>

      <QueryForm />
      <br />

      <LoadForm to="faculty" from="faculty.csv" />
      <LoadForm to="themes" from="themes.xml" />
      <br />

      <Update label="Drop citations" endpoint="drop/citations" />
      <Update label="Pull citations" endpoint="pull" />
    </div>
  );
}

function QueryForm({ _ }) {
  // const [v, setV] = React.useState(3);
  return (
    <form
      onSubmit={async function (event) {
        event.preventDefault();
        const request = {};
        [...event.target]
          .filter((input) => input.value !== "" && input.name !== "")
          .forEach((input) => {
            request[input.name] = input.value.split(",").map((w) => w.trim());
          });
        const result = await fetch(
          `/query?q=${encodeURIComponent(JSON.stringify(request))}`,
          {
            method: "GET",
          }
        ).then((response) => response.json());
        console.log(result);
      }}
    >
      <QueryText name="faculty" label="surnames" />
      <br />
      <QueryText name="keyword" label="keywords" />
      <br />

      <input type="submit" value="Search" />
    </form>
  );
}

function QueryText({ name, label }) {
  return (
    <label>
      {label}: <input name={name} type="text" />
    </label>
  );
}

function LoadForm({ to, from }) {
  return (
    <form
      method="post"
      action={`/update/load/${to}`}
      encType="multipart/form-data"
    >
      <label>
        upload {from}: <input name="file" type="file" />
      </label>
      <br />
      <label>
        (password) <input name="key" type="text" />
      </label>
      <input type="submit" value="Upload" />
    </form>
  );
}

function Update({ endpoint, label }) {
  return (
    <form action={`/update/${endpoint}`}>
      <Password />
      <input type="submit" value={label} />
    </form>
  );
}

function Password({ _ }) {
  return (
    <label>
      (password) <input name="key" type="text" />
    </label>
  );
}

export default App;
