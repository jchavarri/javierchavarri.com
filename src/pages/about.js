import React from "react";
import { graphql } from "gatsby";
import Layout from "../components/layout";

export default ({ data }) => (
  <Layout>
    <h1>
      <span role="img" aria-label="Ghost">
        👻
      </span>
    </h1>
    <p>
      I'm a software engineer building products with{" "}
      <a
        href="http://reasonml.github.io/"
        target="_blank"
        rel="noopener noreferrer"
      >
        ReasonML
      </a>
      , currently at{" "}
      <a href="http://ahrefs.com/" target="_blank" rel="noopener noreferrer">
        Ahrefs
      </a>
      . Before that, I worked on frontend infrastructure at{" "}
      <a href="http://webflow.com/" target="_blank" rel="noopener noreferrer">
        Webflow
      </a>
      .
    </p>
    <p>
      I am the author of some{" "}
      <a
        href="https://github.com/jchavarri"
        target="_blank"
        rel="noopener noreferrer"
      >
        open source projects
      </a>{" "}
      (mostly in Reason) and I am also a contributor to UI or tooling projects
      like{" "}
      <a
        href="https://github.com/revery-ui/revery"
        target="_blank"
        rel="noopener noreferrer"
      >
        Revery
      </a>{" "}
      or{" "}
      <a
        href="https://github.com/cristianoc/genType"
        target="_blank"
        rel="noopener noreferrer"
      >
        genType
      </a>
      .
    </p>
    <p>
      I am very interested on programming languages and language design from the
      purely technical perspective, but also and especially from the human point
      of view: how they affect and relate with remote work, team cooperation,
      products maintenance and evolution, software quality, and most importantly
      personal happiness.
    </p>
    <p>
      If you want to chat, you can find me on{" "}
      <a
        href="https://twitter.com/javierwchavarri/"
        target="_blank"
        rel="noopener noreferrer"
      >
        Twitter
      </a>
      ,{" "}
      <a
        href="https://www.linkedin.com/in/javier-chávarri-alvarez-91739150"
        target="_blank"
        rel="noopener noreferrer"
      >
        LinkedIn
      </a>
      , or email (javier dot chavarri at gmail).
    </p>
  </Layout>
);

export const query = graphql`
  query {
    site {
      siteMetadata {
        title
      }
    }
  }
`;
