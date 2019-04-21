import Typography from "typography";
import anneTheme from "typography-theme-st-annes";

anneTheme.baseFontSize = "18px";
anneTheme.overrideThemeStyles = ({ rhythm }, options) => ({
  // "h1,h2,h3": {
  //   marginTop: rhythm(1),
  // },
  // "p": {
  //   marginBottom: rhythm(2),
  // },
  a: {
    color: "#333",
    textDecoration: "underline",
  },
});
const typography = new Typography(anneTheme);

export default typography;
export const rhythm = typography.rhythm;
