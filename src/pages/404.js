import React from "react";
import Layout from "../components/layout";
import gif from "../../static/media/404.gif";

export default ({ data }) => (
  <Layout>
    <h1>
      <span role="img" aria-label="Ghost">
        🤔 Not found
      </span>
    </h1>
    <p>
      It seems this page does not exist. Go back to the{" "}
      <a href="https://javierchavarri.com">home page</a> or check{" "}
      <a
        href="https://giphy.com/search/random"
        target="_blank"
        rel="noopener noreferrer"
      >
        these other random gifs
      </a>
      <br />
    </p>
    <img src={gif} alt="404" style={{ width: "100%" }} />
  </Layout>
);
