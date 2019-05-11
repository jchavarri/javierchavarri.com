import React from "react";
import Helmet from "react-helmet";
import { graphql } from "gatsby";
import Layout from "../components/layout";

export default ({ data }) => {
  const post = data.markdownRemark;
  return (
    <>
      <Helmet title={post.frontmatter.title}>
        {/* Twitter Card tags */}
        <meta name="twitter:title" content={post.frontmatter.title} />
        <meta name="twitter:description" content={post.frontmatter.subtitle} />
        <meta name="twitter:image" content={post.frontmatter.imghero} />
      </Helmet>
      <Layout>
        <div>
          <h1>{post.frontmatter.title}</h1>
          {post.frontmatter.subtitle && <h2>{post.frontmatter.subtitle}</h2>}
          <p>
            <em>{post.frontmatter.date}</em>
          </p>
          <div dangerouslySetInnerHTML={{ __html: post.html }} />
        </div>
      </Layout>
    </>
  );
};

export const query = graphql`
  query($slug: String!) {
    markdownRemark(fields: { slug: { eq: $slug } }) {
      html
      frontmatter {
        title
        subtitle
        date(formatString: "MMMM DD, YYYY")
      }
    }
  }
`;
