export const formatEffectValueHelper = (value: any): string => {
  if (typeof value === "number") {
    // Display whole numbers without decimal, others as is (or to a fixed precision if needed)
    return value % 1 === 0 ? value.toFixed(0) : value.toString();
  }
  return String(value);
};

export const formatDescription = (
  desc: string | undefined,
  effects: Record<string, any> | undefined,
): string => {
  if (!desc) return "";
  let formattedDesc = desc;

  // 1. Replace placeholders like @EffectName@ or @EffectName*100@
  if (effects) {
    const placeholderRegex = /@([^@]+)@/g;
    let resultString = "";
    let lastIndex = 0;
    let match;

    while ((match = placeholderRegex.exec(formattedDesc)) !== null) {
      resultString += formattedDesc.substring(lastIndex, match.index); // Append text before placeholder

      const placeholderContent = match[1]; 
      if (placeholderContent.startsWith("TFTUnitProperty")) {
        // Skip this placeholder entirely - don't add anything to resultString
        lastIndex = placeholderRegex.lastIndex;
        continue;
      }
      
      let replacementValue = match[0]; // Default to original placeholder if not processed

      if (placeholderContent.endsWith("*100")) {
        const key = placeholderContent.slice(0, -4); // Remove '*100'
        // Case insensitive lookup
        const effectKey = Object.keys(effects).find(k => k.toLowerCase() === key.toLowerCase());
        if (effectKey !== undefined && typeof effects[effectKey] === "number") {
          replacementValue = (effects[effectKey] * 100).toFixed(0);
        }
      } else {
        const key = placeholderContent;
        // Case insensitive lookup
        const effectKey = Object.keys(effects).find(k => k.toLowerCase() === key.toLowerCase());
        if (effectKey !== undefined) {
          replacementValue = formatEffectValueHelper(effects[effectKey]);
        }
      }
      resultString += replacementValue;
      lastIndex = placeholderRegex.lastIndex;
    }
    resultString += formattedDesc.substring(lastIndex); // Append remaining text
    formattedDesc = resultString;
  }

  // 2. Handle specific TFT tags and other HTML-like tags

  // Styling for specific keywords/concepts
  formattedDesc = formattedDesc.replace(
    /<TFTKeyword>(.*?)<\/TFTKeyword>/gi,
    '<span class="text-yellow-500 font-medium">$1</span>',
  );
  formattedDesc = formattedDesc.replace(
    /<tftbold>(.*?)<\/tftbold>/gi,
    '<em class="font-semibold">$1</em>', // Italic + Semi-bold
  );
  formattedDesc = formattedDesc.replace(
    /<magicDamage>(.*?)<\/magicDamage>/gi,
    '<span class="text-purple-400">$1</span>',
  );
  formattedDesc = formattedDesc.replace(
    /<physicalDamage>(.*?)<\/physicalDamage>/gi,
    '<span class="text-red-500">$1</span>',
  );
  formattedDesc = formattedDesc.replace(
    /<trueDamage>(.*?)<\/trueDamage>/gi,
    '<span class="text-white font-semibold">$1</span>',
  );
  formattedDesc = formattedDesc.replace(
    /<TFTBonus>(.*?)<\/TFTBonus>/gi,
      // '<span class="text-green-400">$1</span>',
      ""
  );
  formattedDesc = formattedDesc.replace(
    /<TFTShadowItemPenalty>(.*?)<\/TFTShadowItemPenalty>/gi,
    '<span class="text-red-400 italic">$1</span>', 
  );
  formattedDesc = formattedDesc.replace(
    /<healing>(.*?)<\/healing>/gi,
    '<span class="text-green-300">$1</span>',
  );
  formattedDesc = formattedDesc.replace(
    /<TFTHighlight>(.*?)<\/TFTHighlight>/gi,
    '<span class="text-blue-300 font-medium">$1</span>',
  );

  // Basic handling for <li> items (convert to bullet points on new lines)
  formattedDesc = formattedDesc.replace(/<li>(.*?)<\/li>/gi, "<br />&bull; $1"); 

  // Tags to strip while keeping their content.
  const tagsToStripKeepContent = [
    "tftitemrules",
    "rules",
    "TFTRadiantItemBonus",
    "TFTShadowItemBonus",
    "TFTPassive",
    "scaleLevel",
    "scaleShimmer",
    "TFTTrackerLabel",
    "spellPassive",
    "spellActive",
    "scaleHealth",
    "keyword",
    "tftrules",
    "active",
    "ShowIfNot\\.TFT14_Mob_IsActive_T3",
    "ShowIf\\.TFT14_Mob_IsActive_T3",
    "ShowIfNot\\.TFTUnitProperty\\.Item:TFT10_BlingActive",
    "ShowIf\\.TFTUnitProperty\\.Item:TFT10_BlingActive",
    "ShowIfCustom\\.Set=TFTSetEventCT",
  ];

  tagsToStripKeepContent.forEach((tag) => {
    const contentKeepingRegex = new RegExp(`<${tag}[^>]*>(.*?)</${tag}>`, "gi");
    formattedDesc = formattedDesc.replace(contentKeepingRegex, "$1");
    const stripTagOnlyRegex = new RegExp(`</?${tag}[^>]*>`, "gi");
    formattedDesc = formattedDesc.replace(stripTagOnlyRegex, "");
  });

  // Consolidate multiple <br> tags (and variants) into one
  formattedDesc = formattedDesc.replace(/(<br\s*\/?>\s*)+/gi, "<br />");
  
  // Replace &nbsp; with regular space
  formattedDesc = formattedDesc.replace(/&nbsp;/g, " ");
  
  // Clean up any remaining problematic HTML entities or malformed tags
//   formattedDesc = formattedDesc.replace("<br />", "\n");
  formattedDesc = formattedDesc.replace(/&lt;\/br&gt;/gi, "");

  return formattedDesc;
};