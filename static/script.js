ReactDOM.render(
  <App message="Gillings Search Tool" />,
  document.getElementById("root")
);

function App(props) {
  return (
    <div>
      <h1>{props.message}</h1>

      <p>
        <form action="/query">
          <label for="faculty">surname</label>
          <input name="faculty" type="text" value="Agans" />
          <br />
          <label for="keyword">keyword</label>
          <input name="keyword" type="text" value="Reliability" />
          <br />
          <input type="submit" value="Search" />
        </form>
        <br />
      </p>

      <p>
        <form
          method="post"
          action="/update/load/faculty"
          enctype="multipart/form-data"
        >
          <label for="file">upload faculty.csv: </label>
          <input name="file" type="file" />
          <br />
          <label for="key">(password)</label>
          <input name="key" type="text" />
          <input type="submit" value="Upload" />
        </form>
        <form
          method="post"
          action="/update/load/themes"
          enctype="multipart/form-data"
        >
          <label for="file">upload themes.xml: </label>
          <input name="file" type="file" />
          <br />
          <label for="key">(password)</label>
          <input name="key" type="text" />
          <input type="submit" value="Upload" />
        </form>
      </p>

      <p>
        <form action="/update/drop/citations">
          <input type="submit" value="Drop citations" />
        </form>
        <form action="/update/pull">
          <input type="submit" value="Pull citations" />
        </form>
      </p>
    </div>
  );
}
