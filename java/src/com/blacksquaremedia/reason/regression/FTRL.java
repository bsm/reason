package com.blacksquaremedia.reason.regression;

import java.lang.Math;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Collections;
import com.blacksquaremedia.reason.CoreProtos;
import com.blacksquaremedia.reason.core.Example;
import com.google.protobuf.ProtocolStringList;
import net.jpountz.xxhash.XXHash64;
import net.jpountz.xxhash.XXHashFactory;

public class FTRL {
	private FTRLProtos.Optimizer optimizer;
	private XXHashFactory hasher;
	private CoreProtos.Model model;
	private ArrayList<String> predictors;
	private int[] offsets;
	private double alpha;
	private double beta;
	private double l1;
	private double l2;

	public FTRL(FTRLProtos.Optimizer optimizer) {
		this(optimizer, 0.1, 1.0, 1.0, 0.1);
	}

	public FTRL(FTRLProtos.Optimizer optimizer, double alpha, double beta, double l1, double l2) {
		this.optimizer = optimizer;
		this.hasher = XXHashFactory.fastestInstance();
		this.model = this.optimizer.getModel();
		this.predictors = new ArrayList<String>(model.getFeaturesCount()-1);
		this.offsets = new int[model.getFeaturesCount()-1];
		this.alpha = alpha;
		this.beta = beta;
		this.l1 = l1;
		this.l2 = l2;

		// Build a sorted list of feature names
		String targetName = optimizer.getTarget();
		for (String name : this.model.getFeaturesMap().keySet()) {
			if (name != targetName) {
				this.predictors.add(name);
			}
		}
		Collections.sort(this.predictors);

		// Pre-calculate predictor offsets
		int pos = 0;
		for (int i = 0; i < this.predictors.size(); i++) {
			this.offsets[i] = pos;

      CoreProtos.Feature feature = this.model.getFeaturesOrThrow(this.predictors.get(i));
      switch (feature.getKind()) {
      case CATEGORICAL:
        pos += (feature.getHashBuckets() + feature.getVocabularyCount());
        break;
      case NUMERICAL:
        pos += 1;
        break;
      }
		}
	}

	// Predict returns a single probability for this example.
	public double predict(Example example) {
		double wTx = 0.0;

		// Iterate over sorted predictors.
		for (int i = 0; i < this.predictors.size(); i++) {
			wTx += this.increment(example, this.predictors.get(i), this.offsets[i]);
		}

		// Apply sigmoid function to calculate probability.
		return 1.0 / (1.0 + Math.exp(-Math.max(Math.min(wTx, 35), -35)));
	}

	private double increment(Example x, String predictor, int offset) {
    CoreProtos.Feature feature = this.model.getFeaturesOrThrow(predictor);

    // Get the example value for the predictor
    Object obj = x.getExampleValue(predictor);
    if (obj == null) {
      return 0.0;
    }

    // Default bucket offset and value
    int bucket = offset;
    double value = 1.0;

    // Adjust bucket offset and value
    switch (feature.getKind()) {
      case CATEGORICAL:
      	bucket += this.categoricalOffset(feature, obj.toString());
        break;
      case NUMERICAL:
		    if (obj instanceof Number) {
		      value = ((Number) obj).doubleValue();
		    }
        break;
      default:
      	return 0.0;
    }

		double weight = this.optimizer.getWeights(bucket);
		double sign = weight < 0 ? -1.0 : 1.0;
		double abs = weight * sign;
		if (abs <= this.l1) {
			return 0.0;
		}

		double sum = Math.sqrt(this.optimizer.getSums(bucket));
		double step = this.l2 + (this.beta+sum)/this.alpha;
		return sign * (this.l1 - abs) / step * value;
	}


	private int categoricalOffset(CoreProtos.Feature feature, String value) {
    ProtocolStringList vocabulary = feature.getVocabularyList();

    // Check if that string is included in the vocabulary; use the index if so.
    int index = vocabulary.indexOf(value);
    if (index > -1) {
      return index;
    }

    // If the string isn't stored in the vocabulary, use bucket hashing.
    int numBuckets = feature.getHashBuckets();
    byte[] data = value.getBytes();
    XXHash64 hash = this.hasher.hash64();
    long hashed = hash.hash(data, 0, data.length, 0);
    BigInteger bigInt = new BigInteger(Long.toUnsignedString(hashed));
    BigInteger bucket = bigInt.remainder(BigInteger.valueOf(numBuckets));

    return bucket.intValue() + vocabulary.size();
	}
}
