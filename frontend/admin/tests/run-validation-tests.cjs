require("esbuild-register/dist/node").register({ target: "es2022" });

const { validateFeatureForm } = require("../src/lib/validation/feature");

const assert = (condition, message) => {
  if (!condition) {
    throw new Error(message);
  }
};

const buildFeature = () => ({
  name: "checkout-button",
  description: "Toggle the new checkout button.",
  active: true,
  variants: [
    { name: "control", weight: 50 },
    { name: "experiment", weight: 50 },
  ],
});

const run = () => {
  const valid = validateFeatureForm(buildFeature());
  assert(valid.isValid, "Valid feature should pass");
  assert(!valid.errors.name, "Valid feature should not flag name");

  const missingBasics = {
    name: "",
    description: "",
    active: true,
    variants: [{ name: "", weight: "" }],
  };
  const basicsResult = validateFeatureForm(missingBasics);
  assert(!basicsResult.isValid, "Missing basics should fail");
  assert(
    basicsResult.errors.variants[0].weight === "Weight is required.",
    "Missing weight should be reported",
  );

  const duplicate = buildFeature();
  duplicate.variants = [
    { name: "Control", weight: 50 },
    { name: "control", weight: 50 },
  ];
  const duplicateResult = validateFeatureForm(duplicate);
  assert(!duplicateResult.isValid, "Duplicate names should fail");
  assert(
    duplicateResult.errors.variants[0].name ===
      "Each variant needs a unique name.",
    "Duplicate name message missing",
  );

  const badTotal = buildFeature();
  badTotal.variants = [
    { name: "control", weight: 30 },
    { name: "experiment", weight: 30 },
  ];
  const badTotalResult = validateFeatureForm(badTotal);
  assert(!badTotalResult.isValid, "Total weight mismatch should fail");
  assert(
    badTotalResult.errors.totalWeight?.includes("currently 60") ?? false,
    "Total weight error should mention current sum",
  );

  const longDescription = buildFeature();
  longDescription.description = "x".repeat(520);
  const descriptionResult = validateFeatureForm(longDescription);
  assert(!descriptionResult.isValid, "Long description should fail");
  assert(
    descriptionResult.errors.description?.includes("under 500 characters") ??
      false,
    "Description length message missing",
  );
};

run();

console.log("All validation checks passed.");
