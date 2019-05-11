import React from "react";
import { Helmet } from "react-helmet";
import { css } from "@emotion/core";
import { Link, graphql } from "gatsby";
import { rhythm } from "../utils/typography";
import Layout from "../components/layout";
import { formatReadingTime } from "../utils/helpers";

export default ({ data }) => {
  const {
    site: {
      siteMetadata: { siteUrl, description, title },
    },
  } = data;
  return (
    <>
      <Helmet defaultTitle={title} titleTemplate={`%s | ${title}`}>
        <meta property="og:title" content={title} />
        <meta property="og:description" content={description} />
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
          <p>{description}</p>
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
                <em>{`${node.frontmatter.date} • ${formatReadingTime(
                  node.timeToRead
                )}`}</em>
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
        description
      }
    }
    allMarkdownRemark(sort: { fields: [frontmatter___date], order: DESC }) {
      totalCount
      edges {
        node {
          id
          timeToRead
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
