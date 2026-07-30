import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:shimmer/shimmer.dart';

class LazyImage extends StatelessWidget {
  final String? imageUrl;
  final double width;
  final double height;
  final BoxFit fit;
  final double borderRadius;
  final String? placeholderAsset;

  const LazyImage({
    super.key,
    this.imageUrl,
    this.width = double.infinity,
    this.height = 200,
    this.fit = BoxFit.cover,
    this.borderRadius = 8,
    this.placeholderAsset,
  });

  @override
  Widget build(BuildContext context) {
    if (imageUrl == null || imageUrl!.isEmpty) {
      return _placeholder();
    }
    return ClipRRect(
      borderRadius: BorderRadius.circular(borderRadius),
      child: CachedNetworkImage(
        imageUrl: imageUrl!,
        width: width,
        height: height,
        fit: fit,
        placeholder: (context, url) => _shimmer(),
        errorWidget: (context, url, error) => _placeholder(),
        fadeInDuration: const Duration(milliseconds: 300),
        fadeOutDuration: const Duration(milliseconds: 100),
        memCacheWidth: width.toInt() > 0 ? width.toInt() : null,
        memCacheHeight: height.toInt() > 0 ? height.toInt() : null,
      ),
    );
  }

  Widget _shimmer() {
    return Shimmer.fromColors(
      baseColor: Colors.grey[300]!,
      highlightColor: Colors.grey[100]!,
      child: Container(width: width, height: height, color: Colors.white),
    );
  }

  Widget _placeholder() {
    return Container(
      width: width,
      height: height,
      decoration: BoxDecoration(
        color: Colors.grey[200],
        borderRadius: BorderRadius.circular(borderRadius),
      ),
      child: placeholderAsset != null
          ? Image.asset(placeholderAsset!, fit: BoxFit.contain)
          : const Icon(Icons.image, size: 48, color: Colors.grey),
    );
  }
}
