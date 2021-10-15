ReactDOM.render(<App message="Hello!" />, document.getElementById("root"));

function App(props) {
  return (
    <div>
      <h1>{props.message}</h1>
      <form action="http://localhost:3001/query/">
        <input name="faculty" type="text" />
        <input type="submit" value="Search" />
      </form>
    </div>
  );
}
