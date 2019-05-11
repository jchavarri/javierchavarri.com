import React from "react";
import { Helmet } from "react-helmet";
import { css } from "@emotion/core";
import { Link, graphql } from "gatsby";
import { rhythm } from "../utils/typography";
import Layout from "../components/layout";

export default ({ data }) => {
  const {
    site: {
      siteMetadata: { siteUrl, title },
    }
  } = data;
  return (
    <>
      <Helmet defaultTitle={title} titleTemplate={`%s | ${title}`}>
        <meta property="og:url" content={siteUrl} />
        <meta property="og:type" content="website" />
      </Helmet>

      <Layout>
        <div>
          <p>
            Hi there!{" "}
            <span role="img" aria-label="Wave">
              👋
            </span>{" "}
          </p>
          <p>
            Welcome my site, where I write about the challenges of creating
            software products, UI engineering, type systems and programming
            languages, especially JavaScript and Reason.
          </p>
          {data.allMarkdownRemark.edges.map(({ node }) => (
            <div key={node.id}>
              <Link
                to={node.fields.slug}
                css={css`
                  color: inherit;
                  text-decoration: none;
                `}
              >
                <h3
                  css={css`
                    margin-bottom: ${rhythm(1 / 4)};
                  `}
                >
                  {node.frontmatter.title}
                </h3>
              </Link>
              <p
                css={css`
                  color: #bbb;
                  margin-bottom: ${rhythm(1 / 4)};
                `}
              >
                <em>{node.frontmatter.date}</em>
              </p>
              <p>{node.excerpt}</p>
            </div>
          ))}
        </div>
      </Layout>
    </>
  );
};

export const query = graphql`
  query {
    site {
      siteMetadata {
        siteUrl
        title
      }
    }
    allMarkdownRemark(sort: { fields: [frontmatter___date], order: DESC }) {
      totalCount
      edges {
        node {
          id
          frontmatter {
            title
            date(formatString: "MMMM DD, YYYY")
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
