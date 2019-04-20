import React from "react";
import { Helmet } from "react-helmet";
import { css } from "@emotion/core";
import { StaticQuery, Link, graphql } from "gatsby";
import { rhythm } from "../utils/typography";
import "./layout.css";

export default ({ children }) => (
  <StaticQuery
    query={graphql`
      query {
        site {
          siteMetadata {
            title
          }
        }
      }
    `}
    render={data => (
      <div
        css={css`
          margin: 0 auto;
          max-width: 42rem;
          padding: ${rhythm(1)};
          padding-top: ${rhythm(1.5)};
        `}
      >
        <Helmet>
          <html lang="en" />
          <meta charSet="utf-8" />
          <title>{data.site.siteMetadata.title}</title>

          {/* Favicon stuff from https://favicon.io/favicon-generator/ */}
          <link
            rel="apple-touch-icon"
            sizes="180x180"
            href="/apple-touch-icon.png"
          />
          <link
            rel="icon"
            type="image/png"
            sizes="32x32"
            href="/favicon-32x32.png"
          />
          <link
            rel="icon"
            type="image/png"
            sizes="16x16"
            href="/favicon-16x16.png"
          />
          <link rel="manifest" href="/site.webmanifest" />
        </Helmet>
        <Link to={`/`}>
          <h3
            css={css`
              margin-bottom: ${rhythm(2)};
              display: inline-block;
              font-style: normal;
            `}
          >
            {data.site.siteMetadata.title}
          </h3>
        </Link>
        <Link
          to={`/about/`}
          css={css`
            float: right;
          `}
        >
          <h3>About</h3>
        </Link>
        {children}
      </div>
    )}
  />
);
