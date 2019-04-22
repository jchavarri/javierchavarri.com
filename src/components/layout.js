import React from "react";
import { css } from "@emotion/core";
import { StaticQuery, Link, graphql } from "gatsby";
import { rhythm } from "../utils/typography";
import SiteMetadata from "./site-metadata";
import "./layout.css";

export default ({ children, pathname }) => (
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
        <SiteMetadata pathname={pathname} />
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
