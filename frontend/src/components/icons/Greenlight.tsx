// Official Greenlight emblem from Blockstream's Blueprint design system:
// https://design.blockstream.com/styleguide/branding/core-lightning#greenlight
// Paths copied from greenlight_emblem_on_transparent_rgb.svg in the official
// CoreLightningGreenlightLightningChargeAssets.zip. The primary greenlight
// yellow (#F0D003) and white are used on a dark badge so both halves remain
// visible on Alby Hub's light setup surface.
export const GreenlightIcon = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 2000 2000"
      fill="none"
      aria-hidden="true"
    >
      <rect width="2000" height="2000" fill="#121217" />
      <path
        fill="#F0D003"
        d="M1078.4,656l216.8-256.7L380,1091.5h102l0,0h370l646.9-742.7L1078.4,656z M434.2,1073.3l741.5-560.9l-159.7,189.2L507,1073.3H434.2z"
      />
      <path
        fill="#FFFFFF"
        d="M1620,908.5h-330.7l0,0h-141.6l-646.6,742.7l420.3-306.9l-216.8,256.7L1620,908.5z M1493,926.7h72.9l-741.8,560.9l159.7-189.2L1493,926.7z"
      />
    </svg>
  );
};
