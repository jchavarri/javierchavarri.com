import React from "react";
import { css } from "@emotion/core";
import { Link, graphql } from "gatsby";
import { rhythm } from "../utils/typography";
import Layout from "../components/layout";

export default ({ data }) => {
  return (
    <Layout>
      <div>
        <p>
          Hi there!{" "}
          <span role="img" aria-label="Wave">
            👋
          </span>{" "}
        </p>
        <p>
          Welcome to this site, where I write (i.e. rant) about software
          products, UI engineering, and programming languages —mainly JavaScript
          and Reason.
        </p>
        {data.allMarkdownRemark.edges.map(({ node }) => (
          <div key={node.id}>
            <Link
              to={node.fields.slug}
              css={css`
                color: inherit;
              `}
            >
              <h3
                css={css`
                  margin-bottom: ${rhythm(1 / 4)};
                `}
              >
                {node.frontmatter.title}
              </h3>
              <p
                css={css`
                  color: #bbb;
                  margin-bottom: ${rhythm(1 / 4)};
                `}
              >
                <em>{node.frontmatter.date}</em>
              </p>
              <p>{node.excerpt}</p>
            </Link>
          </div>
        ))}
      </div>
    </Layout>
  );
};

export const query = graphql`
  query {
    allMarkdownRemark(sort: { fields: [frontmatter___date], order: DESC }) {
      totalCount
      edges {
        node {
          id
          frontmatter {
            title
            date(formatString: "DD MMMM, YYYY")
          }
          fields {
            slug
          }
          excerpt
        }
      }
    }
  }
`;
